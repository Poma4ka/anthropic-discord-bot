import { CompletionMessage } from './completion-message';

export type GetPreviousMessage = () => Promise<CompletionMessage | null>;
