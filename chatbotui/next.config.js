/** @type {import('next').NextConfig} */
const nextConfig = {
    reactStrictMode: true,
};

module.exports = {
    nextConfig,
    env: {
        SERVER_URL: 'https://chatbot-tele-server.herokuapp.com',
    },
};
