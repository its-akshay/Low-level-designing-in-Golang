## Entity Identification Framework

When given any LLD problem, extract entities using this mental model:

```
Nouns     → Structs / Models
Verbs     → Methods / Interfaces
Adjectives → Fields / Attributes
Ownership → Composition vs Aggregation
Behavior  → Interface contracts
```

"A parking lot has multiple floors. Each floor has parking spots.
Vehicles can park and unpark. Tickets are issued on entry."

Nouns    → ParkingLot, Floor, Spot, Vehicle, Ticket, EntryGate, ExitGate
Verbs    → Park(), Unpark(), IssueTicket(), CalculateFee()
Adjectives → Available, Occupied, Compact, Large
Ownership → ParkingLot HAS-MANY Floors → Composition
Behavior  → ParkingStrategy interface for spot allocation


Writing maintainable backend systems:

1. Define contracts first (interfaces before implementations)
2. Keep structs focused (Single Responsibility)
3. Inject dependencies (testability, flexibility)
4. Prefer composition over embedding for flexibility
5. Use unexported fields with exported methods (encapsulation)
6. Return errors explicitly — never panic in library code
7. Use context.Context for cancellation and timeouts
8. Write table-driven tests alongside implementation


| OOP Concept | Java/C++ | Go Equivalent |
| --- | --- | --- |
| Class | `class Foo {}` | `type Foo struct {}` |
| Private field | `private int x` | `x int` (lowercase) |
| Public field | `public int X` | `X int` (uppercase) |
| Constructor | `new Foo()` | `NewFoo()` factory function |
| Inheritance | `extends BaseClass` | Embedding `BaseStruct` |
| Interface | `implements IFoo` | Implicit satisfaction |
| Abstract class | `abstract class` | Interface + base struct |
| Polymorphism | Method overriding | Interface implementation |