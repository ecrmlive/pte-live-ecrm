import COS from 'cos-js-sdk-v5';

import {
  applyLiveVodUploadApi,
  commitLiveVodUploadApi,
  type LiveVodApplyUploadResult,
} from '#/api/core/live-vod';

import { captureVideoFirstFrame } from './video-cover';

function mediaTypeFromFile(file: File) {
  const name = file.name.toLowerCase();
  const ext = name.includes('.') ? name.split('.').pop() : '';
  if (ext === 'm3u8') return 'm3u8';
  if (ext === 'mov') return 'mov';
  if (ext === 'webm') return 'webm';
  return 'mp4';
}

function putObjectToCos(
  body: Blob | File,
  storagePath: string,
  info: LiveVodApplyUploadResult,
  cred: NonNullable<LiveVodApplyUploadResult['temp_certificate']>,
  onProgress?: (progressData: { percent?: number }) => void,
) {
  return new Promise<void>((resolve, reject) => {
    const cos = new COS({
      getAuthorization(_options, callback) {
        callback({
          ExpiredTime: cred.expired_time,
          SecurityToken: cred.token,
          StartTime: Math.floor(Date.now() / 1000),
          TmpSecretId: cred.secret_id,
          TmpSecretKey: cred.secret_key,
        });
      },
    });
    const key = String(storagePath || '').replace(/^\//, '');
    if (!key) {
      reject(new Error('缺少上传路径'));
      return;
    }
    cos.putObject(
      {
        Body: body,
        Bucket: info.storage_bucket!,
        Key: key,
        Region: info.storage_region!,
        onProgress(progressData) {
          onProgress?.(progressData);
        },
      },
      (err) => {
        if (err) reject(err);
        else resolve();
      },
    );
  });
}

export interface VodUploadProgress {
  fileName: string;
  percent: number;
  phase: '' | 'commit' | 'cover' | 'done' | 'preparing' | 'uploading';
}

export async function uploadVodVideo(
  file: File,
  options: {
    media_name?: string;
    media_type?: string;
    onProgress?: (payload: VodUploadProgress) => void;
  } = {},
) {
  if (!file) {
    throw new Error('请选择视频文件');
  }
  const report = (percent: number, phase: VodUploadProgress['phase']) => {
    options.onProgress?.({
      fileName: file.name,
      percent: Math.min(100, Math.max(0, Math.round(percent))),
      phase,
    });
  };

  report(0, 'preparing');
  let coverBlob: Blob | null = null;
  try {
    coverBlob = await captureVideoFirstFrame(file);
  } catch {
    coverBlob = null;
  }
  report(5, 'preparing');

  const applyPayload: Record<string, string> = {
    media_name: options.media_name || file.name,
    media_type: options.media_type || mediaTypeFromFile(file),
  };
  if (coverBlob) {
    applyPayload.cover_type = 'Jpg';
  }
  const info = await applyLiveVodUploadApi(applyPayload);
  const cred = info.temp_certificate || {};
  if (!info.vod_session_key || !info.storage_bucket || !cred.secret_id) {
    throw new Error('申请云点播上传失败');
  }
  report(8, 'uploading');

  await putObjectToCos(file, info.media_storage_path!, info, cred, (progressData) => {
    const ratio = Number(progressData?.percent) || 0;
    report(8 + ratio * 82, 'uploading');
  });

  if (coverBlob && info.cover_storage_path) {
    report(92, 'cover');
    await putObjectToCos(coverBlob, info.cover_storage_path, info, cred);
  }

  report(96, 'commit');
  const result = await commitLiveVodUploadApi({
    vod_session_key: info.vod_session_key,
  });
  if (!result.file_id || !result.media_url) {
    throw new Error('确认云点播上传失败');
  }
  report(100, 'done');
  return result;
}
