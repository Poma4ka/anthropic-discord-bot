import { MessageRoleEnum } from '../enum';

export interface CompletionAttachment {
  content: Buffer;
  contentType?: string;
  name?: string;
}

export interface CompletionMessage {
  content: string;
  attachments?: CompletionAttachment[];
  role: MessageRoleEnum;
}
