import Input from "@mui/joy/Input";

const Search = () => {
  return (
    <Input
      placeholder="Search for cuisine or a dish"
      sx={{
        fontSize: "14px",
        height: "40px",
        width: "500px",
        "--Input-focusedThickness": "0px",
      }}
    />
  );
};

export default Search;
