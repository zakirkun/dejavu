import Link from 'next/link'
import { Rocket, Zap, Shield, BarChart } from 'lucide-react'

export default function Home() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-purple-50">
      {/* Header */}
      <header className="container mx-auto px-4 py-6">
        <nav className="flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <Rocket className="h-8 w-8 text-blue-600" />
            <span className="text-2xl font-bold text-gray-900">Dejavu</span>
          </div>
          <div className="flex items-center space-x-4">
            <Link href="/login" className="text-gray-600 hover:text-gray-900">
              Login
            </Link>
            <Link
              href="/register"
              className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
            >
              Get Started
            </Link>
          </div>
        </nav>
      </header>

      {/* Hero Section */}
      <section className="container mx-auto px-4 py-20 text-center">
        <h1 className="text-6xl font-bold text-gray-900 mb-6">
          Deploy Instantly,
          <br />
          <span className="text-blue-600">Scale Effortlessly</span>
        </h1>
        <p className="text-xl text-gray-600 mb-8 max-w-2xl mx-auto">
          Platform deployment modern yang memudahkan Anda deploy aplikasi Next.js, Node.js, Go, dan lainnya dalam hitungan detik.
        </p>
        <div className="flex items-center justify-center space-x-4">
          <Link
            href="/register"
            className="px-8 py-4 bg-blue-600 text-white rounded-lg hover:bg-blue-700 text-lg font-semibold"
          >
            Start Deploying Free
          </Link>
          <Link
            href="#features"
            className="px-8 py-4 border-2 border-gray-300 text-gray-700 rounded-lg hover:border-gray-400 text-lg font-semibold"
          >
            Learn More
          </Link>
        </div>
      </section>

      {/* Features */}
      <section id="features" className="container mx-auto px-4 py-20">
        <h2 className="text-4xl font-bold text-center text-gray-900 mb-12">
          Fitur Unggulan
        </h2>
        <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-8">
          <FeatureCard
            icon={<Zap className="h-10 w-10 text-blue-600" />}
            title="Deploy Cepat"
            description="Deploy aplikasi Anda dalam hitungan detik dengan auto-detection framework"
          />
          <FeatureCard
            icon={<Rocket className="h-10 w-10 text-blue-600" />}
            title="Auto Scaling"
            description="Aplikasi Anda akan auto-scale berdasarkan traffic tanpa konfigurasi rumit"
          />
          <FeatureCard
            icon={<Shield className="h-10 w-10 text-blue-600" />}
            title="SSL Otomatis"
            description="HTTPS gratis untuk semua deployment dengan Let's Encrypt"
          />
          <FeatureCard
            icon={<BarChart className="h-10 w-10 text-blue-600" />}
            title="Real-time Logs"
            description="Monitor deployment dan aplikasi Anda secara real-time"
          />
        </div>
      </section>

      {/* Supported Frameworks */}
      <section className="container mx-auto px-4 py-20 bg-white rounded-2xl shadow-lg my-20">
        <h2 className="text-4xl font-bold text-center text-gray-900 mb-12">
          Framework yang Didukung
        </h2>
        <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-8 items-center justify-items-center">
          <FrameworkLogo name="Next.js" />
          <FrameworkLogo name="Node.js" />
          <FrameworkLogo name="Bun" />
          <FrameworkLogo name="Go" />
          <FrameworkLogo name="PHP" />
          <FrameworkLogo name="Static" />
        </div>
      </section>

      {/* Pricing */}
      <section className="container mx-auto px-4 py-20">
        <h2 className="text-4xl font-bold text-center text-gray-900 mb-12">
          Harga Transparan
        </h2>
        <div className="grid md:grid-cols-3 gap-8 max-w-5xl mx-auto">
          <PricingCard
            name="Starter"
            price="Gratis"
            features={[
              '3 Projects',
              '100 GB Bandwidth',
              'SSL Otomatis',
              'Community Support',
            ]}
          />
          <PricingCard
            name="Pro"
            price="$19/bulan"
            features={[
              'Unlimited Projects',
              '1 TB Bandwidth',
              'SSL Otomatis',
              'Priority Support',
              'Custom Domains',
            ]}
            highlighted
          />
          <PricingCard
            name="Enterprise"
            price="Custom"
            features={[
              'Everything in Pro',
              'Unlimited Bandwidth',
              'Dedicated Support',
              'SLA 99.99%',
              'Custom Infrastructure',
            ]}
          />
        </div>
      </section>

      {/* CTA */}
      <section className="container mx-auto px-4 py-20 text-center">
        <div className="bg-blue-600 rounded-2xl p-12 text-white">
          <h2 className="text-4xl font-bold mb-4">
            Siap untuk Deploy?
          </h2>
          <p className="text-xl mb-8 opacity-90">
            Mulai deploy aplikasi Anda sekarang, gratis!
          </p>
          <Link
            href="/register"
            className="inline-block px-8 py-4 bg-white text-blue-600 rounded-lg hover:bg-gray-100 text-lg font-semibold"
          >
            Get Started Now
          </Link>
        </div>
      </section>

      {/* Footer */}
      <footer className="container mx-auto px-4 py-8 border-t border-gray-200 mt-20">
        <div className="text-center text-gray-600">
          <p>&copy; 2025 Dejavu. All rights reserved.</p>
        </div>
      </footer>
    </div>
  )
}

function FeatureCard({ icon, title, description }: { icon: React.ReactNode, title: string, description: string }) {
  return (
    <div className="p-6 bg-white rounded-xl shadow-md hover:shadow-lg transition-shadow">
      <div className="mb-4">{icon}</div>
      <h3 className="text-xl font-semibold text-gray-900 mb-2">{title}</h3>
      <p className="text-gray-600">{description}</p>
    </div>
  )
}

function FrameworkLogo({ name }: { name: string }) {
  return (
    <div className="text-center">
      <div className="w-20 h-20 bg-gray-100 rounded-lg flex items-center justify-center mb-2">
        <span className="text-2xl font-bold text-gray-700">{name[0]}</span>
      </div>
      <p className="text-sm text-gray-600">{name}</p>
    </div>
  )
}

function PricingCard({ name, price, features, highlighted = false }: { name: string, price: string, features: string[], highlighted?: boolean }) {
  return (
    <div className={`p-8 rounded-xl ${highlighted ? 'bg-blue-600 text-white shadow-2xl scale-105' : 'bg-white shadow-md'}`}>
      <h3 className={`text-2xl font-bold mb-2 ${highlighted ? 'text-white' : 'text-gray-900'}`}>{name}</h3>
      <p className={`text-4xl font-bold mb-6 ${highlighted ? 'text-white' : 'text-blue-600'}`}>{price}</p>
      <ul className="space-y-3 mb-8">
        {features.map((feature, i) => (
          <li key={i} className="flex items-center">
            <svg className={`w-5 h-5 mr-2 ${highlighted ? 'text-white' : 'text-green-500'}`} fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
            </svg>
            {feature}
          </li>
        ))}
      </ul>
      <Link
        href="/register"
        className={`block w-full py-3 text-center rounded-lg font-semibold ${
          highlighted
            ? 'bg-white text-blue-600 hover:bg-gray-100'
            : 'bg-blue-600 text-white hover:bg-blue-700'
        }`}
      >
        Get Started
      </Link>
    </div>
  )
}

