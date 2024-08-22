import * as types from "@common/types.ts"
import { Mutex } from "./utils";

export function exhausted(d: never) {
    console.log(d)
    throw new Error("unreachable: " + JSON.stringify(d));
}

export type Message = {
    Typ: types.Varient.ExeNotFound,
    Val: types.TExeNotFound,
} | {
    Typ: types.Varient.UserLogin,
    Val: types.TUserLogin,
} | {
    Typ: types.Varient.Err,
    Val: types.TErr,
} | {
    Typ: types.Varient.Unknown,
    Val: unknown,
}

// type ValType<T extends types.Varient> = Message extends { Typ: T; Val: infer V } ? V : never;
type ValType<T extends types.Varient> = 
    T extends types.Varient.ExeNotFound ? types.TExeNotFound :
    T extends types.Varient.UserLogin ? types.TUserLogin :
    T extends types.Varient.Err ? types.TErr :
    unknown;

type Callback<T extends types.Varient> = (res: ValType<T>) => PromiseLike<void>;
type DisableCallback = () => Promise<void>;

export class Server {
    ws: WebSocket;

    protected wait: Promise<void>;
    protected constructor() {
        this.ws = new WebSocket(`ws://localhost:${6200}/${"ws"}`);

        this.ws.addEventListener('message', async (e) => {
            let mesg: Message = JSON.parse(e.data);
            // @ts-ignore
            mesg.Val = JSON.parse(mesg.Val);

            await this.handle_message(mesg);
        });

        let resolve: () => {};
        this.wait = new Promise(r => {
            resolve = r as () => {};
        });
        this.ws.addEventListener('open', async (_e) => {
            resolve();
        });
    }

    static async new() {
        let self = new Server();
        await self.wait;
        return self;
    }

    // @ts-ignore
    callbacks = new Map<Message["Typ"], [number, Callback<Message["Typ"]>][]>();
    callback_mutex: Mutex = new Mutex();
    callback_id = 1;
    async add_callback<T extends Message["Typ"]>(type: T, cb: Callback<typeof type>): Promise<DisableCallback> {
        let id = await this.callback_mutex.runExclusive(async () => {
            let id = this.callback_id++;

            let callbacks = this.callbacks.get(type) ?? null;
            if (callbacks == null) {
                this.callbacks.set(type, callbacks);
            } else {
                callbacks.push([id, cb]);
            }

            return id;
        });

        return async () => {
            await this.callback_mutex.runExclusive(async () => {
                let callbacks = this.callbacks.get(type) ?? [];
                callbacks = callbacks.filter(e => e[0] != id);
                this.callbacks.set(type, callbacks);
            });
        };
    }
    async handle_message(msg: Message) {
        console.log(msg);

        let callbacks = this.callbacks.get(msg.Typ) ?? [];
        for (let callback of callbacks) {
            await callback[1](msg.Val);
        }
        if (callbacks.length > 0) {
            return;
        }

        switch (msg.Typ) {
            case types.Varient.ExeNotFound:
            case types.Varient.Err:
                break;
            case types.Varient.Unknown:
            case types.Varient.UserLogin: {
                throw new Error(`message type '${msg.Typ}' can't be handled here`);
            } break;
            default:
                throw exhausted(msg);
        }
    }

    send_message(msg: Message) {
        console.log(msg);
        msg.Val = JSON.stringify(msg.Val);
        let json = JSON.stringify(msg);
        this.ws.send(json);
    }
}

// @ts-ignore
export let server: Server = null;
export async function init() {
    server = await Server.new();
}

