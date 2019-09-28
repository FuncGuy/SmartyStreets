First Series:

1)  Create an envelope with address [street1] and send it over input channel.
2)  Handler recevies envelope from input channel and assigns to application verifier and also writes to output channel.

Note:-

-> as of now application verifier is not doing anything except storing the envelope.


Second Series:

1)  Previous tomato is all about “Handler” listening the message[Enevelope] from input channel and passing it over “Verifier” and to Output channel.

2)  Now “Verifier” comes into an action here once the message comes to the “Verifier” it has to construct an HTTP request and send it to the “SmartyStreets” API for validations.

3)  Here an “GET” request is created by appending the message parameters to query string.


Third Series:

1) This tomato is about mocking out smartystreets API a fake client FakeHttpClient is being built it accepts a request and produces valid json response.

2) Response has been constructed and tested for all the scenarios.


Fourth Series:

1) This series is about refractoring the test cases.
2) And decorating the request by adding authentication info.

Fitfth Series: (Cami and Mike)

1) It is about finishing up the “Authorizer”
2) And starts to work on the “Sequencer”
3) And the logic of sequenceer is envlopes can come in any order for exmple 3,4,0,2,1
    the job of sequencer is reorder[0,1,2,3,4] the envelopes in sequence[contigous]

Sixth Series:

1) Jonathan returns from his meeting and mike quickly walks him through the Sequencer.
2) Then they start working on the Writer.
3) Notice the design of Writer
4) And refractroing the Test cases.
5) They seem to have no need for a mocking framework as they write their own thats intresting.
6) And finally they start to work on the Reader.
7) Read from csv file and write to envelope 
	
Seventh Series:

1) All about Reader

Eighth Series:

1) Final Cleanup of EOF troubles that we had in the last episode.
2) Writing an integration test it helps us to WIRE up the system.



Go learnings:

1)  When a method has pointer(*) value for example:

	func (wallet *Wallet) TestHello(t *testing.T) {
                          // code
} 

It means -> it will update the state of wallet [attributes]

2) When a method returns an address(&) reference for example:

        	func (wallet *Wallet) TestHello(t *testing.T) {
                       return &NewWallet()
}  
 
It means  -> it is returning that object /reference. 

3) “defer” is like finally in java it will execute once the surrounding function has been executed.

4) creating map : make(map[int]*Envelope) or make(map[int]Envelope{})
