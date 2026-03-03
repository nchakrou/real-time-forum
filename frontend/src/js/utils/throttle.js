export function throttle(fn, delay) {
  let timer;
  return function (...args) {
    if (timer) {
      return;
    }
    fn(...args);
    timer = setTimeout(() => {
      timer = null;
    }, delay);
  };
}
