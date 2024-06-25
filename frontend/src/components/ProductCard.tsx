import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

import { Button } from "@/components/ui/button";

type Product = {
  id: number;
  title: string;
  category: string;
};

export function ProductCard(props: Product) {
  const { title, category } = props;
  return (
    <Card className="w-[250px] h-[250px] flex flex-col justify-evenly items-center">
      <CardHeader>
        <CardTitle>{title}</CardTitle>
      </CardHeader>
      <CardContent>
        <CardDescription>{category}</CardDescription>
      </CardContent>
      <CardFooter>
        <Button variant="default">Add to cart</Button>
      </CardFooter>
    </Card>
  );
}
