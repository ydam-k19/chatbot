import Layout from '../../components/Layout/Layout';
import { useRouter } from 'next/router';
import { useState, useEffect } from 'react';

export default function Vietjet() {
    const router = useRouter();
    const [isMobile, setIsMobile] = useState(false);

    let { meetingID, chatID, messageID, name } = router.query;
    chatID = Number(chatID);
    messageID = Number(messageID);
    const [reason, setReason] = useState();

    useEffect(() => {
        let details = navigator.userAgent;
        let regexp = /android|iphone|kindle|ipad/i;
        setIsMobile(regexp.test(details));
    }, []);

    const handleSubmit = async (e) => {
        e.preventDefault();
        const res = await fetch(`${process.env.SERVER_URL}/meeting/refuse-task`, {
            method: 'POST',
            body: JSON.stringify({ meeting_id: meetingID, chat_id: chatID, message_id: messageID, name, reason }),
            headers: {
                'Content-Type': 'application/json',
            },
        });
        if (isMobile) {
            window.close();
            return;
        }
        if (res.status === 200) {
            return router.push('/success');
        }
    };
    const handleChange = (e) => {
        setReason(e.target.value);
    };

    return (
        <Layout>
            <form onSubmit={handleSubmit} method="post">
                <input
                    type="text"
                    name="reason"
                    value={reason}
                    onChange={handleChange}
                    placeholder="Enter the reason:"
                />
                <button type="submit">Submit</button>
            </form>
        </Layout>
    );
}
