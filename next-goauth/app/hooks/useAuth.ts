import { useState, useEffect } from "react";
import { User, getCurrentUser } from "../actions/actions";

export default function useAuth() {
  const [user, setUser] = useState<User>();

  useEffect(() => {
    getCurrentUser().then((res) => setUser(res));
  }, []);

  return {
    user,
  };
}
