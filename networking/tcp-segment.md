# tcp-segment

TCP 协议的 PDU 称为数据段(Segment)，即 TCP 协议将需要发送的数据，放在一个一个 Segment 里面，从发送端发往接收端。

TCP 在中文语境下有一个不是问题的问题，即 "粘包"。

首先 TCP 是一个基于流的协议。TCP 可以把一条很长的数据拆成多个 Segment 进行发送，也可以把三条数据汇集成两个 Segment 进行发送，他让应用层开发者从这些繁冗的细节中解放了出来。但至于怎么把数据拆成 Segment，又是怎么发送数据，都不是一个应用层的开发者应该关心的。但是如果一个应用层开发者，从应用层的逻辑(服务端发一条数据，客户端端收一条数据)来理解 TCP 的时候，就会出现一些偏差。

这个可怜的开发者会发现，为什么我调用了两次 Send，结果只要一次 Recv 就收完了？我的"包"被吃了吗？

显然不是

通过 Send 发送出去的数据一般称为 Message，而 TCP 的数据单位被称为 Segment。显然 Mesaage 和 Segment 可以是多对多的关系。例如

```toml
Message_A = [A1, A2, A3, A4]

Segments = [
  [A1, A2, A3],
  [A4],
]
```

也可能是

```toml
Message_A = [A1, A2, A3, A4]
Message_B = [B1, B2, B3, B4, B5, B6, B7]
Message_C = [C1, C2, C3]

Segments = [
  [A1, A2, A3, A4, B1, B2, B3],
  [B4, B5, B6, B7, C1, C2, C3],
]
```

有人就把 Message_B 随着 Message_A 一起发出去的行为，称为 "粘包"。

所以，"粘包" 是一个问题吗？

显然不是。这是个正常的 TCP 行为，问题出把 Message 和 Segment 理解成一对一的关系。

这个错误的理解会导致什么问题呢？最常见的情况是，我设计了一种基于 TCP 的协议。为了处理这个协议进行网络编程，当我 Recv 了一个 Segment 之后，我读到了协议的 Message_A，但是 Message_B 是粘着 Message_A 一起来的，如果读完 Message_A 之后我认为任务完成了，这样一来 Message_B 的头部分就被丢掉了。所以导致的最直接的问题，发生在了 TCP 解(自定义协议的)包(的Message) 的上。

这其实变成了一个协议设计的问题，为了使 decoder 更加可靠/容易编写，就有了这么几种设计

- fixed length decoder
协议每次发送固定长度的 Message，如果需要发送的 Segment 长度不够时用空字符弥补，并且保证不超过缓冲区，接受方每次按固定长度区接受数据

适合长度多数情况下固定的数据，如传感器的监测数据，但是在传输不定长数据时会有较大带宽的浪费

- delimiter based decoder
在 body 部分主动塞入 Message 的分隔符，用来标记 Message 的边界

缺点是 Message 本体中的分隔符需要转义

- length field decoder
在 header 部分塞入包的长度，解包时根据 header 的信息区分不同 Message

适用于在发送时就知道长度的 Message

TODO: Nagle's algorithm
