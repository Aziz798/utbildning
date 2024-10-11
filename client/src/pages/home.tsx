import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import img from "@/assets/img.jpg";
import logo from "@/assets/logo.svg";
import { useState } from "react";
import { Eye, EyeOff } from "lucide-react";
import axiosInstance from "@/lib/axios-instance";
import { AxiosError } from "axios";
import { useNavigate } from "react-router-dom";
import { ModeToggle } from "@/components/shared/mode-toggle";
export default function Home() {
    const [data, setData] = useState({
        email: "",
        password: "",
    });
    const [errors, setErrors] = useState<null | { [key: string]: string }>(null);
    const [showPassword, setShowPassword] = useState(false);
    const nav = useNavigate();
    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        setErrors(null);
        e.preventDefault();
        try {
            const res = await axiosInstance.post("login", data);
            setData({
                email: "",
                password: "",
            })
            localStorage.setItem("accessToken", res.data.access_token);
            localStorage.setItem("refreshToken", res.data.refresh_token);
            nav("/dashboard");

        } catch (error) {
            console.log(error)
            if (error instanceof AxiosError) {
                setErrors(error.response?.data);
                console.log(errors);
            }
        }
    };
    
    return (
        <div className="grid w-full min-h-screen grid-cols-1 lg:grid-cols-2">
            <div className="bg-[#dbe8f0] block lg:hidden">
                <img src={img} alt="ethiopian woman smiling" width={1920}
                    height={1080}
                    className="h-full w-full object-cover"
                    style={{ aspectRatio: "1920/1080", objectFit: "cover" }}
                />
            </div>
            <div className="flex items-center justify-center py-12">
                <div className="mx-auto grid w-[350px] gap-6">
                    <div className="flex items-center justify-between gap-5 w-full">
                    <img src={logo} alt="" />
                    <ModeToggle/>
                    </div>
                    <div className="grid gap-2 text-center">
                        <h1 className="text-3xl font-bold">Login</h1>
                        <p className="text-balance text-muted-foreground">
                            Enter your email below to login to your account
                        </p>
                        {
                            errors ?
                                <p className="text-balance text-red-500">Invalid email or password</p> : null
                        }
                    </div>
                    <form onSubmit={handleSubmit}>
                        <div className="grid gap-4">
                            <div className="grid gap-2">
                                <Label htmlFor="email">Email</Label>
                                <Input
                                    id="email"
                                    type="email"
                                    placeholder="m@example.com"
                                    required
                                    value={data.email}
                                    onChange={(e) =>
                                        setData({ ...data, email: e.target.value })
                                    }
                                />
                            </div>
                            <div className="grid gap-2">
                                <div className="flex items-center">
                                    <Label htmlFor="password">Password</Label>
                                </div>
                                <div className="flex flex-row items-center gap-1">
                                    <Input
                                        id="password"
                                        type={showPassword ? "text" : "password"}
                                        required
                                        value={data.password}
                                        onChange={(e) => setData({ ...data, password: e.target.value })}
                                    />
                                    <Button
                                        variant="outline"
                                        size="sm"
                                        className=""
                                        type="button"
                                        onClick={() => setShowPassword((x) => !x)}
                                    >
                                        {showPassword ? <EyeOff /> : <Eye />}
                                    </Button>
                                </div>
                            </div>
                            <Button type="submit" className="w-full">
                                Login
                            </Button>
                        </div>
                    </form>
                </div>
            </div>
            <div className="hidden lg:block lg:bg-[#dbe8f0]">
                <img
                    src={img}
                    alt="Login Image"
                    width={1920}
                    height={1080}
                    className="h-full w-full object-cover"
                    style={{ aspectRatio: "1920/1080", objectFit: "cover" }}
                />
            </div>
        </div>
    )
}