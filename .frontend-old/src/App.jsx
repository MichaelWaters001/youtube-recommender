import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Home from "./pages/Home";
import CreatorDetails from "./pages/CreatorDetails";
import Navbar from "./components/Navbar";

function App() {
  return (
    <Router>
      <Navbar />
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/creators/:id" element={<CreatorDetails />} />
      </Routes>
    </Router>
  );
}

export default App;