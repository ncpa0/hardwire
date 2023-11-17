import { If, condition } from "../../src/components/gotmpl-generator/if";
import { Range } from "../../src/components/gotmpl-generator/range";
import { DynamicFragment } from "../../src/components/stream";

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

const Fallback = () => {
  return <div class="loading-indicator">Loading...</div>;
};

export const StreamedFragment = async () => {
  return (
    <DynamicFragment<User>
      require="User"
      fallback={<Fallback />}
      render={(user) => {
        return (
          <div>
            <h3>{user.username.toString()}</h3>
            <If
              condition={condition((b) =>
                b.and(
                  user.isPremium,
                  b.or(
                    b.greaterThan(user.age, b.value(18)),
                    b.lessEqualThan(user.sinceSignedUp, b.value(30))
                  )
                )
              )}
            >
              <Range
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
    ></DynamicFragment>
  );
};
