
// - [Simple TypeScript Mutex Implementation - DEV Community](https://dev.to/0916dhkim/simple-typescript-mutex-implementation-5544)
export class Mutex {
    private _queue: {
        resolve: (release: ReleaseFunction) => void;
    }[] = [];

    private _isLocked = false;

    acquire() {
        return new Promise<ReleaseFunction>((resolve) => {
            this._queue.push({ resolve });
            this._dispatch();
        });
    }

    async runExclusive<T>(callback: () => Promise<T>) {
        const release = await this.acquire();
        try {
            return await callback();
        } finally {
            release();
        }
    }

    private _dispatch() {
        if (this._isLocked) {
            return;
        }
        const nextEntry = this._queue.shift();
        if (!nextEntry) {
            return;
        }
        this._isLocked = true;
        nextEntry.resolve(this._buildRelease());
    }

    private _buildRelease(): ReleaseFunction {
        return () => {
            this._isLocked = false;
            this._dispatch();
        };
    }
}

type ReleaseFunction = () => void;
