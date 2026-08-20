const Card = () => {
  return (
      <div className="w-full max-w-2xl mt-8 sm:mt-12 p-5 sm:p-8 bg-white rounded-2xl shadow-xl border border-neutral-200">
        <h2 className="text-xl font-semibold text-center text-neutral-800">
          Paste the URL to be shortened
        </h2>

        <div className="mt-6 flex flex-col sm:flex-row gap-3">
          <input
            type="url"
            placeholder="https://example.com/very-long-url"
            className="w-full flex-1 px-5 py-4 text-neutral-700 border border-neutral-300 rounded-xl outline-none focus:ring-2 focus:ring-neutral-900 focus:border-transparent transition-all"
          />

          <button className="w-full sm:w-auto px-6 py-4 bg-neutral-900 text-white font-medium rounded-xl hover:bg-neutral-800 transition-colors cursor-pointer">
            Shorten URL
          </button>
        </div>

        <p className="mt-4 text-sm text-neutral-400 text-center">
          Your shortened URL will appear here.
        </p>
      </div>
  )
}

export default Card
