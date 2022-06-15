import Layout from '../../components/Layout/Layout';
import { useRouter } from 'next/router';
import { useState, useEffect } from 'react';

export default function RegisterAdmin() {
    const router = useRouter();
    const { chatID, name } = router.query;
    const [code, setCode] = useState();
    const [isMobile, setIsMobile] = useState(false);
    const [msgErr, setMsgErr] = useState({});

    useEffect(() => {
        let details = navigator.userAgent;
        let regexp = /android|iphone|kindle|ipad/i;
        setIsMobile(regexp.test(details));
        setMsgErr(document.getElementsByClassName('msg-err')[0]);
    }, []);

    const handleSubmit = async (e) => {
        e.preventDefault();
        const res = await fetch(`${process.env.SERVER_URL}/register/admin`, {
            method: 'POST',
            body: JSON.stringify({ chat_id: Number(chatID), name, code }),
            headers: {
                'Content-Type': 'application/json',
            },
        });

        if (isMobile) {
            window.close();
            return;
        }
        if (res.status === 201) {
            return router.push('/success');
        } else {
            msgErr.innerHTML = 'Code is invalid >.< Try again!';
        }
    };
    const handleChange = (e) => {
        msgErr.innerHTML = '';
        setCode(e.target.value);
    };

    return (
        <Layout>
            <div>
                <form onSubmit={handleSubmit} method="post">
                    <input
                        type="text"
                        name="code"
                        value={code}
                        onChange={handleChange}
                        placeholder="Enter private code:"
                    />
                    <button type="submit">Submit</button>
                </form>
                <div className="msg-err"></div>
            </div>
        </Layout>
    );
}
