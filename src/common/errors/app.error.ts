export class AppError extends Error {
  public readonly message: string;

  constructor(message: string) {
    super();
    this.message = message;
  }

  getStack(error: Error = this) {
    if (error === this) {
      if (!this.stack) {
        Error.captureStackTrace(error, this.constructor);
      }
      return this.stack;
    }

    if (error.stack) {
      this.stack = error.stack;
    }

    return this.stack;
  }
}
