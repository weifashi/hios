/**
 * 页面专用
 */
window.systemInfo = window.systemInfo || {}
const webTs = {
    /**
     * 接口地址
     * @param str
     * @returns {string|string|*}
     */
    apiUrl(str) {
        if (
            str.substring(0, 2) === "//" ||
            str.substring(0, 7) === "http://" ||
            str.substring(0, 8) === "https://" ||
            str.substring(0, 6) === "ftp://" ||
            str.substring(0, 1) === "/"
        ) {
            return str
        }
        if (typeof window.systemInfo.apiUrl === "string") {
            str = window.systemInfo.apiUrl + str
        } else {
            str = window.location.origin + "/api/" + str
        }
        while (str.indexOf("/../") !== -1) {
            str = str.replace(/\/(((?!\/).)*)\/\.\.\//, "/")
        }
        return str
    },
}
export default webTs
