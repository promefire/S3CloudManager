import { inject } from 'vue';

export function useNotification() {
  const notify = inject('notify');
  return {
    success: (msg, duration) => notify(msg, 'success', duration),
    error: (msg, duration) => notify(msg, 'error', duration),
    info: (msg, duration) => notify(msg, 'info', duration),
    warning: (msg, duration) => notify(msg, 'warning', duration)
  };
}
