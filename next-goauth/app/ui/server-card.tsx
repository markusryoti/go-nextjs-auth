import { getCurrentUser } from "../actions/actions";

export default async function ServerCard() {
  const data = await getCurrentUser();

  return (
    <div>
      <h2 className="text-4xl">User inside server component</h2>
      <pre>{JSON.stringify(data)}</pre>
    </div>
  );
}
