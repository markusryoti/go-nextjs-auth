import ClientCard from "../ui/client-card";
import LogoutButton from "../ui/logout-button";
import ServerCard from "../ui/server-card";

export default async function Dashboard() {
  return (
    <div className="flex flex-col gap-8">
      <h1 className="text-5xl mb-8">Dashboard</h1>
      <div className="flex flex-col gap-8">
        <ServerCard />
        <ClientCard />
        <LogoutButton />
      </div>
    </div>
  );
}
