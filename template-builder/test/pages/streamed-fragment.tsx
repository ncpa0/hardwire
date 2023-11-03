import { If, condition } from "../../src/components/gotmpl-generator/if";
import { Map } from "../../src/components/gotmpl-generator/map";
import { Stream } from "../../src/components/stream";

type User = {
  username: string;
  age: number;
  sinceSignedUp: number;
  isPremium: boolean;
  premiumLinks: Array<{
    name: string;
    url: string;
  }>;
};

export const StreamedFragment = async () => {
  return (
    <Stream<User>
      require="User"
      render={(user) => {
        return (
          <div>
            <h3>{user.username.toString()}</h3>
            <If
              condition={condition((b) =>
                b.and(
                  user.isPremium,
                  b.or(
                    b.gt(user.age, b.value(18)),
                    b.le(user.sinceSignedUp, b.value(30))
                  )
                )
              )}
            >
              <Map
                data={user.premiumLinks}
                render={(link) => {
                  return (
                    <a href={link.url.toString()}>{link.name.toString()}</a>
                  );
                }}
              />
            </If>
          </div>
        );
      }}
    ></Stream>
  );
};
