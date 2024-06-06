"use client";

import useAuth from "../hooks/useAuth";

export default function ClientCard() {
  const { user } = useAuth();

  return (
    <div>
      <h2 className="text-4xl">User inside client component</h2>
      <pre>{JSON.stringify(user)}</pre>
    </div>
  );
}
