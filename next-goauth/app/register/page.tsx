import { register } from "../actions/actions";

export default function Register() {
  return (
    <div>
      <h1 className="text-5xl mb-8">Register</h1>
      <form action={register} className="flex flex-col gap-4">
        <div className="">
          <label htmlFor="email">Email</label>
          <input name="email" type="text" className="border h-8 w-96" />
        </div>
        <div>
          <label htmlFor="password">Password</label>
          <input name="password" type="password" className="border h-8 w-96" />
        </div>
        <div>
          <button type="submit" className="bg-slate-500 text-white p-4">
            Submit
          </button>
        </div>
      </form>
    </div>
  );
}
