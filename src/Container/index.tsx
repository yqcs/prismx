import React from "react";

const Container: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  return <div style={{ width: 1080, margin: "0 auto" }}>{children}</div>;
};

export default Container;
