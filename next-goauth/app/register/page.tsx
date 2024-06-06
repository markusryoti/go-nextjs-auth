import { register } from "../actions/actions";

export default function Register() {
  return (
    <div>
      <h1 className="text-5xl mb-8">Register</h1>
      <form action={register} className="flex flex-col gap-4">
        <div className="flex flex-col">
          <label htmlFor="email" className="text-gray-700 px-4 py-2">
            Email
          </label>
          <input
            name="email"
            type="email"
            className="border w-96 p-4 h-12 rounded-xl"
          />
        </div>
        <div className="flex flex-col">
          <label htmlFor="password" className="text-gray-700 px-4 py-2">
            Password
          </label>
          <input
            name="password"
            type="password"
            className="border w-96 p-4 h-12 rounded-xl"
          />
        </div>
        <div>
          <button
            type="submit"
            className="bg-blue-500 text-white p-4 rounded-xl"
          >
            Submit
          </button>
        </div>
      </form>
    </div>
  );
}
