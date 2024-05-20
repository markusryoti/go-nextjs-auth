"use server";

import jwt from "jsonwebtoken";

const secret = "verysecret";

type Claims = {
  session: string;
  exp: number;
};

export async function verifyToken(rawToken: string) {
  const payload = jwt.verify(rawToken, secret);
  const { session, exp } = payload as Claims;

  const expiresAt = exp * 1000;

  return { session, expiresAt };
}

export async function signToken(session: string) {
  const token = jwt.sign({ session }, secret);
  return token;
}
