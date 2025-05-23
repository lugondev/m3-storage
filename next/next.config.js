/** @type {import('next').NextConfig} */
const nextConfig = {
	images: {
		remotePatterns: [
			{
				protocol: 'https',
				hostname: 'picsum.photos',
				port: '',
				pathname: '/**',
			},
		]
	},
	poweredByHeader: false,
	reactStrictMode: true,
};

module.exports = nextConfig;
