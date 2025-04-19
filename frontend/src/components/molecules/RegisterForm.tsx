import { FormEvent, useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { useAppDispatch, useAppSelector } from "@/lib/hooks/AppHooks";

import { cn } from "@/lib/helpers/utils";
import { Button } from "@/components/atoms/shadCN/button";
import { Input } from "@/components/atoms/shadCN/input";
import { Label } from "@/components/atoms/shadCN/label";
import { registerThunk } from "@/services/auth/thunks";

export function RegisterForm({
  className,
  ...props
}: React.ComponentPropsWithoutRef<"form">) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [err, setErr] = useState("");

  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const token = useAppSelector((s) => s.auth.token);

  useEffect(() => {
    if (token) navigate("/");
  }, [token, navigate]);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    const res = await dispatch(registerThunk({ email, password }));
    if (registerThunk.rejected.match(res)) {
      setErr(res.payload as string);
    }
  };

  return (
    <form
      onSubmit={handleSubmit}
      className={cn("flex flex-col gap-6", className)}
      {...props}
    >
      <div className="flex flex-col items-center gap-2 text-center">
        <h1 className="text-2xl font-bold">Create an account</h1>
        <p className="text-balance text-sm text-muted-foreground">
          Enter your email below to create your account
        </p>
      </div>
      <div className="grid gap-6">
        <div className="grid gap-2">
          <Label htmlFor="email">Email</Label>
          <Input
            id="email"
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="m@example.com"
            required
          />
        </div>
        <div className="grid gap-2">
          <Label htmlFor="password">Password</Label>
          <Input
            id="password"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        {err && <p className="text-sm text-red-600">{err}</p>}
        <Button type="submit" className="w-full">
          Create account
        </Button>
      </div>

      <div className="text-center text-sm">
        Already have an account?{" "}
        <a href="/login" className="underline underline-offset-4">
          Login
        </a>
      </div>
    </form>
  );
}
