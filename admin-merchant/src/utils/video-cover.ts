/** 从本地视频文件截取首帧，返回 JPEG Blob（用于云点播封面上传）。 */
export function captureVideoFirstFrame(file: File): Promise<Blob> {
  return new Promise((resolve, reject) => {
    if (!file) {
      reject(new Error('缺少视频文件'));
      return;
    }
    const url = URL.createObjectURL(file);
    const video = document.createElement('video');
    video.preload = 'auto';
    video.muted = true;
    video.playsInline = true;

    const cleanup = () => {
      video.removeAttribute('src');
      video.load();
      URL.revokeObjectURL(url);
    };

    video.onloadeddata = () => {
      video.currentTime = Math.min(0.1, (video.duration || 1) * 0.01);
    };

    video.onseeked = () => {
      try {
        const w = video.videoWidth;
        const h = video.videoHeight;
        if (!w || !h) {
          cleanup();
          reject(new Error('无法读取视频尺寸'));
          return;
        }
        const canvas = document.createElement('canvas');
        canvas.width = w;
        canvas.height = h;
        const ctx = canvas.getContext('2d');
        if (!ctx) {
          cleanup();
          reject(new Error('无法创建画布'));
          return;
        }
        ctx.drawImage(video, 0, 0, w, h);
        canvas.toBlob(
          (blob) => {
            cleanup();
            if (blob) resolve(blob);
            else reject(new Error('无法生成封面'));
          },
          'image/jpeg',
          0.85,
        );
      } catch (err) {
        cleanup();
        reject(err instanceof Error ? err : new Error('截取封面失败'));
      }
    };

    video.onerror = () => {
      cleanup();
      reject(new Error('无法读取视频文件'));
    };

    video.src = url;
  });
}
