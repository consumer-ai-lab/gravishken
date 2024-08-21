import * as types from "@common/types.ts"

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

    async handle_message(msg: Message) {
        console.log(msg);
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

