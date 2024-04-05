import { CompletionMessage, GetPreviousMessage } from '../common';

export type CreateCompletionOptionsDto = {
  getPreviousMessage?: GetPreviousMessage;
  message: CompletionMessage;
  signal?: AbortSignal;
};
export type CreateCompletionResultDto = {
  chunk: string;
};
