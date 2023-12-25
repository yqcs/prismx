import "flex.css/dist/data-flex.css";
import React from "react";
import Apply from "./Apply";
import Banner from "./Banner";
import Brand from "./Brand";
import Features from "./Features";
import Footer from "./Footer";

const Home = () => {
  return (
    <>
      <Banner />
      <Features />
      <Brand />
      <Apply />
      <Footer />
    </>
  );
};

export default Home;
