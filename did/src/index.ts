import { dag, Container, Directory, object, func } from "@dagger.io/dagger";

import { DIDUniversalResolver } from "@quarkid/did-resolver";

const QuarkIDEndpoint = "https://node-ssi.buenosaires.gob.ar";

@object()
export class Did {
  /**
   * Return a json string with the DID document given a DID string
   */
  @func()
  async resolve(did: string): Promise<string> {
    const universalResolver = new DIDUniversalResolver({
      universalResolverURL: QuarkIDEndpoint,
    });

    const didDocument = await universalResolver.resolveDID(did);

    return JSON.stringify(didDocument);
  }
}
