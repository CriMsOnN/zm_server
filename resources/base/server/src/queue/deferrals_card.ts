export const DeferralsCard = (deferrals: Record<string, any>) => {
  return (msg: string) => {
    deferrals.presentCard(
      JSON.stringify({
        type: "AdaptiveCard",
        body: [
          {
            type: "Image",
            url: "",
            horizontalAlignment: "center",
          },
          {
            type: "Container",
            items: [
              {
                type: "TextBlock",
                text: "Welcome to our server",
                weight: "bolder",
                size: "medium",
                horizontalAlignment: "center",
              },
              {
                type: "TextBlock",
                text: msg,
                weight: "bolder",
                size: "medium",
                horizontalAlignment: "center",
              },
            ],
            style: "default",
            bleed: true,
            height: "automatic",
            isVisible: true,
          },
        ],
        $schema: "http://adaptivecards.io/schemas/adaptive-card.json",
        version: "1.3",
      }),
    );
  };
};
