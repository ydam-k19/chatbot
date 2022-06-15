import Layout from '../components/Layout/Layout';

export default function Success() {
    setTimeout(() => {
        window.open('tg://msg');
    }, 1000);
    return (
        <Layout>
            <div className="msg-success">Successfully !!</div>
        </Layout>
    );
}
