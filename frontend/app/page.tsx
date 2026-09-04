import Card from "../components/Card/Card";

const Home = () => {
  return (
    <div className="min-h-screen w-full bg-neutral-50 flex flex-col items-center px-4">
      <h1 className="text-3xl sm:text-5xl font-bold text-center mt-16 text-neutral-900">
        URL Shortener
      </h1>

      <Card />
    </div>
  );
};

export default Home;
