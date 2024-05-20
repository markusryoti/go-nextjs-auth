export default function Login() {
  return (
    <div>
      <h1 className="text-5xl">Login</h1>
      <form action="">
        <div>
          <label htmlFor="email">Email</label>
          <input name="email" type="text" />
        </div>
        <div>
          <label htmlFor="password">Password</label>
          <input name="password" type="password" />
        </div>
      </form>
    </div>
  );
}
