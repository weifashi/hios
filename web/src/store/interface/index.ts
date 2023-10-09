import { User } from "../../api/interface/user";

export interface WsMsg {
    action: number     // 消息类型（如：1、上线；2、下线；3、消息）
    data: any          // 消息内容

    type: string       // 客户端类型（如：user）
    uid: number        // 客户端用户ID（会员ID）
    rid: number        // 客户端序号ID（WebSocket ID）
}

export interface GlobalState {
    baseUrl: string
    baseRoute: string
    isLoading: number
    language: string
    themeName: string
    timer: object
    windowActive: any
    windowScrollY: number
    keyboardType: any
    keyboardHeight: any
    windowTouch: any
}

export interface UserState {
    info: User.Info
}

export interface WsState {
    ws: WebSocket,
    msg: WsMsg,
    uid: number,
    rid: number,
    timeout: any,
    random: string,
    openNum: number,
    listener: object,
}
