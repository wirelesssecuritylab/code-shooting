/**
 * 根据key获取cookie值
 * @param key
 * @returns {string}
 */
export function getCookie(key) {
  let match = document.cookie.match(new RegExp('(^|;|\\s)' + key + '=([^;]+)'));
  if (!match) {
    return '';
  }
  return decodeURIComponent(match[2]);
}
