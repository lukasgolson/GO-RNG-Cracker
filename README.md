# Cedar - Random Number Seed Predictor

Cedar (Computationally Efficient Deterministic Algorithm for Random-seed Retrieval) is a powerful tool for predicting random number seeds from sequences generated by Golang's random number generator. With Cedar, you can uncover the seeds used to generate a sequence, which can be useful for predicting future values.

While the idea for Cedar was independently developed, I nearly gave up half way through its development thinking that this feat was impossible. In the midst of deciding to move on to more productive tasks, I did some searching online and stumbled upon the work of @will62794 from 2017.
His blog post [Cracking Go's PRNG For Fun and (Virtual) Profit](https://will62794.github.io/security/hacking/2017/06/30/cracking-golang-prng.html) inspired the approach I took to pre-compute and match sequences.
Inspired by genome sequencing, the innovation that I added to random number cracking is the use of massive Burkhard-Keller trees and a custom binary file format to efficiently search the massive number of possible sequences. 

## How Cedar Works

Cedar employs takes a "smart" brute-force approach to predict random number seeds. It generates all possible sequences from the 2^31 potential seeds available in Golang's random number generator. These sequences are efficiently stored in a BK-tree data structure, allowing for rapid searches. To optimize performance, Cedar utilizes memory-mapped storage and multi-threading.

## Getting Started

To compile Cedar, navigate to the "app/cedar" directory and run the following command:

```bash
go build
```

## Using Cedar

Cedar provides a user-friendly command-line interface (CLI) to interact with the program. Here are some of the available commands:

- **generate:** Generate a sequence of numbers from a given seed.
- **init:** Initialize Cedar by creating the required lookup graphs.
- **search:** Search for the seed that best matches the input sequence.

### Example Usage

Generate a sequence:

```bash
cedar gen -l <length> -s <seed>
```

Initialize Cedar:

```bash
cedar init
```

Search for the best matching seed:

```bash
cedar search <input_sequence>
```

### Command Flags

Cedar also supports various command-line flags to customize its behavior:

- **-c, --cores int:** Specify the number of CPU cores to utilize during processing (default: all of your usable cores).
- **-q, --quiet:** Disable all superficial command text to run Cedar quietly.
- **-t, --toggle:** Display a help message for the toggle functionality.

## Contribution and Development

If you'd like to contribute to Cedar or explore its source code, please refer to our [GitHub repository](https://github.com/lukasgolson/GO-RND-Cracker/).

## License

Cedar is open-source software released under the [MIT License](LICENSE.md).

Cedar is designed to be your go-to tool for predicting random number seeds with ease and efficiency. Feel free to explore its capabilities and contribute to its development. If you have any questions or encounter issues, don't hesitate to reach out. Happy seed cracking!
