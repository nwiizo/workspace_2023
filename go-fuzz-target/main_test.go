package main

import ("math/rand")

func TestEncodeFollowedByDecodeGivesStartingValue(t *testing.T) {
    t.Parallel()
    input := rand.Intn(10)
    encoded := codec.Encode(input)
    t.Logf("encoded value: %#v", encoded)
    want := input
    got := codec.Decode(encoded)
    // after the round trip, we should get what we started with
    if want != got {
        t.Errorf("want %d, got %d", want, got)
    }
}
