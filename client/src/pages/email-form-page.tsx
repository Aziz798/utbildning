
import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from "@/components/ui/card"
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
} from "@/components/ui/alert-dialog"
import axiosInstance from "@/lib/axios-instance"
import { toast } from "sonner"
import { AxiosError } from "axios"

export default function SendEmailForm() {
    const [isDialogOpen, setIsDialogOpen] = useState(false)
    const [subject, setSubject] = useState("")
    const [body, setBody] = useState("")
    const [errors, setErrors] = useState<null | { [key: string]: string }>(null)

    const handleSendEmail = async () => {
        try {
            setIsDialogOpen(false)
            const res = await axiosInstance.post("/send-email-through-json-file", { subject, body })
            console.log(res.data)
            toast("email sent successfully", {
                description: "Your email has been sent successfully",
                duration: 3000,
                action: {
                    label: "Close",

                    onClick: () => console.log("Undo"),
                },
            })

            // Reset form fields after sending
            setSubject("")
            setBody("")
        } catch (error) {
            console.log(error)
            toast("email not sent", {
                description: "Your email has not been sent there were an error",
                action: {
                    label: "Close",
                    onClick: () => console.log("Undo"),
                },
            })
            if (error instanceof AxiosError) {
               setErrors(error.response?.data)
            }
        }
    }

    return (
        <Card className="w-full max-w-2xl m-auto">
            <CardHeader>
                <CardTitle>Compose Email</CardTitle>
                <CardDescription>Write your email message</CardDescription>
            </CardHeader>
            <CardContent>
                <form className="space-y-4">
                    <div className="space-y-2">
                        <Label htmlFor="subject">Subject</Label>
                        <Input
                            id="subject"
                            placeholder="Enter email subject"
                            value={subject}
                            onChange={(e: React.ChangeEvent<HTMLInputElement>) => setSubject(e.target.value)}
                            required
                        />
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="body">Body</Label>
                        <Textarea
                            id="body"
                            placeholder="Write your message here"
                            className="min-h-[200px]"
                            value={body}
                            onChange={(e :React.ChangeEvent<HTMLTextAreaElement>) => setBody(e.target.value)}
                            required
                        />
                    </div>
                </form>
            </CardContent>
            <CardFooter>
                <Button className="w-full" onClick={() => setIsDialogOpen(true)}>Send Email</Button>
            </CardFooter>
        <div>{errors && JSON.stringify(errors)}</div>
            <AlertDialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle>Confirm Send Email</AlertDialogTitle>
                        <AlertDialogDescription>
                            Are you sure you want to send this email? This action cannot be undone.
                        </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                        <AlertDialogCancel>Cancel</AlertDialogCancel>
                        <AlertDialogAction onClick={handleSendEmail}>Send</AlertDialogAction>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>
        </Card>
    )
}
