import { adminPost } from '#/utils/admin-api';

const MessageApi = {
  messageList(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/message/index', data, errorback);
  },
  addMessage(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/message/add', data, errorback);
  },
  editMessage(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/message/edit', data, errorback);
  },
  deleteMessage(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/message/delete', data, errorback);
  },
  fieldList(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/message/field', data, errorback);
  },
  saveField(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/message/saveField', data, errorback);
  },
};

export default MessageApi;
