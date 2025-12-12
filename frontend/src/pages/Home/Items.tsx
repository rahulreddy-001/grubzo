import { Box, Card, Container, Flex, Text } from "@radix-ui/themes";
import { useSelector } from "react-redux";
import type { RootState } from "../../services/store";
import ItemService from "../../services/item/item.service";
import { useEffect } from "react";

function Items() {
  let items = useSelector((s: RootState) => s.item.items);

  useEffect(() => {
    ItemService.getAll();
  }, []);

  let handleAddToCart = (itemID: number, count: number) => {};

  return (
    <Container size="4" style={{ padding: 0 }}>
      <Flex
        direction="row"
        style={{
          height: "100%",
          width: "100%",
        }}
      >
        <Box
          style={{
            background: "red",
            flex: "0 0 20%",
          }}
        >
          Left
        </Box>
        <Box
          style={{
            background: "blue",
            flex: "0 0 80%",
          }}
        >
          {items.map((item) => {
            return (
              <Card>
                <Text>{item.Name}</Text>
                <Text>{item.Description}</Text>
                <Text>{item.Price}</Text>
              </Card>
            );
          })}
          Right
        </Box>
      </Flex>
    </Container>
  );
}

export default Items;
