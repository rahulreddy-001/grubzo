import Button from "@mui/joy/Button";
import { useNavigate } from "react-router";
const plainButtonSX = {
  bgcolor: "transparent",
  color: "grey",
  "&:hover": { bgcolor: "transparent", color: "black" },
  "&:focus-visible": { outline: "none", boxShadow: "none" },
  "--Input-focusedThickness": "0px",
};

const LoginSignUp = () => {
  const navigate = useNavigate();
  return (
    <div>
      <Button
        sx={plainButtonSX}
        onClick={() => {
          navigate("/login");
        }}
      >
        Log in
      </Button>
      <Button
        sx={plainButtonSX}
        onClick={() => {
          navigate("/signup");
        }}
      >
        Sign up
      </Button>
    </div>
  );
};

export default LoginSignUp;
