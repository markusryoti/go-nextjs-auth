"use client";

import { logout } from "../actions/actions";

export default function LogoutButton() {
  const doLogout = async () => {
    await logout();
  };

  return (
    <div>
      <button
        onClick={doLogout}
        className="bg-yellow-400 text-yellow-800 p-4 rounded-xl"
      >
        Logout
      </button>
    </div>
  );
}
