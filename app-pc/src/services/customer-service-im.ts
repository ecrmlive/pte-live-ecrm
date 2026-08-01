import {
  PteLiveIMWebClient,
  type PteChatCredentials,
  type PteChatListener,
  type PteChatMessage,
} from "@pte-live/im-web-sdk";
import {
  fetchCustomerServiceCredential,
  type CustomerServiceIMCredential,
  type CustomerServiceThread,
} from "@/api/customer-service";

function toCredentials(credential: CustomerServiceIMCredential): PteChatCredentials {
  if (
    credential.mode !== "remote" ||
    !credential.sdk_app_id ||
    !credential.im_user_id ||
    !credential.user_sig ||
    !credential.api_url ||
    !credential.ws_url
  ) {
    throw new Error("客服 IM 尚未完成服务端配置，请联系管理员配置 PTE IM 集成");
  }
  return {
    apiUrl: credential.api_url,
    wsUrl: credential.ws_url,
    cosDomain: credential.api_url,
    sdkAppId: credential.sdk_app_id,
    identifier: credential.identifier || credential.im_user_id,
    userId: credential.im_user_id,
    userSig: credential.user_sig,
    expireAt: credential.expire_at,
  };
}

export interface CustomerServiceIMSession {
  thread: CustomerServiceThread;
  conversationID: number;
  client: PteLiveIMWebClient;
  history: () => Promise<PteChatMessage[]>;
  close: () => void;
}

/**
 * 建立商城客服会话。聊天消息只经 pte-live-im-sdk 的加密通道收发；
 * 商城接口只创建线程和签发短期 UserSig。
 */
export async function connectCustomerServiceIM(
  thread: CustomerServiceThread,
  listener: PteChatListener,
): Promise<CustomerServiceIMSession> {
  const getCredential = () => fetchCustomerServiceCredential(thread.thread_id);
  const initialCredential = await getCredential();
  const conversationID = Number(initialCredential.im_conversation_id || thread.im_conversation_id);
  if (!Number.isSafeInteger(conversationID) || conversationID <= 0) {
    throw new Error("客服会话尚未分配 IM 会话，请稍后重试");
  }

  const client = new PteLiveIMWebClient(toCredentials(initialCredential));
  let renewing: Promise<void> | undefined;
  const renewUserSig = async () => {
    if (!renewing) {
      renewing = getCredential()
        .then((credential) => {
          const next = toCredentials(credential);
          client.renewUserSig({ userSig: next.userSig, expireAt: next.expireAt });
        })
        .finally(() => {
          renewing = undefined;
        });
    }
    return renewing;
  };

  const credentialListener: PteChatListener = {
    onUserSigWillExpire: () => void renewUserSig().catch((error: Error) => listener.onError?.(error.message)),
    onUserSigExpired: () => void renewUserSig().catch((error: Error) => listener.onError?.(error.message)),
  };
  client.addListener(listener);
  client.addListener(credentialListener);
  try {
    await client.start();
  } catch (error) {
    client.removeListener(listener);
    client.removeListener(credentialListener);
    client.stop();
    throw error;
  }

  return {
    thread,
    conversationID,
    client,
    history: () => client.history(conversationID),
    close: () => {
      client.removeListener(listener);
      client.removeListener(credentialListener);
      client.stop();
    },
  };
}
