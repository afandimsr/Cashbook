import { jwtDecode } from "jwt-decode";
import type { JwtPayload } from "./JwtPayload";

export function safeDecodeJwt(token?: string): JwtPayload | null {
    if (!token) return null;

    try {
        return jwtDecode<JwtPayload>(token);
    } catch {
        return null;
    }
}
