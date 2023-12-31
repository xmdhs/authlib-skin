export class ApiErr extends Error {
    readonly code: number

    constructor(code: number, msg: string) {
        super(msg)
        this.code = code
    }
}