import { redirect } from "next/navigation";
import { cookies } from "next/headers";
import { signToken } from "../lib/jwt";

async function getUserData() {
  const session = cookies().get("session");
  if (!session) {
    redirect("/login");
  }

  const token = await signToken(session.value);

  const res = await fetch("http://localhost:5000/current-user", {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  if (!res.ok) {
    const body = await res.text();
    console.log(body);
    return undefined;
  }

  const data = await res.json();

  return data;
}

export default async function Dashboard() {
  const data = await getUserData();

  return (
    <div>
      <h1 className="text-5xl mb-8">Dashboard</h1>
      <pre>{JSON.stringify(data)}</pre>
    </div>
  );
}
