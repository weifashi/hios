/**
 * v-longpress
 * 长按指令，长按时触发事件
 */
import utils from "@/utils/utils";
import type { Directive, DirectiveBinding } from "vue";
const isSupportTouch = "ontouchend" in document;
const longpress: Directive = {
    mounted(el, binding: DirectiveBinding) {
        let delay = 500,
            callback = binding.value;
        if (utils.isJson(binding.value)) {
            delay = binding.value.delay || 500;
            callback = binding.value.callback;
        }
        if (typeof callback !== 'function') {
            throw 'callback must be a function'
        }

        // 不支持touch时使用右键
        if (!isSupportTouch) {
            el.__longpressContextmenu__ = (e) => {
                e.preventDefault()
                e.stopPropagation()
                callback(e, el)
            }
            el.addEventListener('contextmenu', el.__longpressContextmenu__);
            return
        }

        // 定义变量
        let pressTimer = null
        let isCall = false
        // 创建计时器（ 500秒后执行函数 ）
        el.__longpressStart__ = (e) => {
            if (e.type === 'click' && e.button !== 0) {
                return
            }
            isCall = false
            if (pressTimer === null) {
                pressTimer = setTimeout(() => {
                    isCall = true
                    callback(e.touches[0], el)
                }, delay)
            }
        }
        // 取消计时器
        el.__longpressCancel__ = (e) => {
            if (pressTimer !== null) {
                clearTimeout(pressTimer)
                pressTimer = null
            }
        }
        // 点击拦截
        el.__longpressClick__ = (e) => {
            if (isCall) {
                e.preventDefault()
                e.stopPropagation()
            }
            el.__longpressCancel__(e)
        }
        // 添加事件监听器
        el.addEventListener('touchstart', el.__longpressStart__)
        // 取消计时器
        el.addEventListener('click', el.__longpressClick__)
        el.addEventListener('touchmove', el.__longpressCancel__)
        el.addEventListener('touchend', el.__longpressCancel__)
        el.addEventListener('touchcancel', el.__longpressCancel__)
    },
    unmounted(el){
        if (!isSupportTouch) {
            el.removeEventListener('contextmenu', el.__longpressContextmenu__)
            delete el.__longpressContextmenu__
            return
        }
        el.removeEventListener('touchstart', el.__longpressStart__)
        el.removeEventListener('click', el.__longpressClick__)
        el.removeEventListener('touchmove', el.__longpressCancel__)
        el.removeEventListener('touchend', el.__longpressCancel__)
        el.removeEventListener('touchcancel', el.__longpressCancel__)
        delete el.__longpressStart__
        delete el.__longpressClick__
        delete el.__longpressCancel__
    }
};

export default longpress;