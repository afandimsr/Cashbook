import { jwtDecode } from "jwt-decode";
import type { JwtPayload } from "./JwtPayload";
import type { TempJwtPayload } from "./TempJwtPayload";

export function safeDecodeJwt(token?: string): JwtPayload | null {
    if (!token) return null;

    try {
        return jwtDecode<JwtPayload>(token);
    } catch {
        return null;
    }
}

export function safeDecodeTempJwt(token?: string): TempJwtPayload | null {
    if (!token) return null;

    try {
        return jwtDecode<TempJwtPayload>(token);
    } catch {
        return null;
    }
}
